# Controller details

Here we will discuss the details of the controller and how it works.

## Controller defined

The controller is defined in the `controller.go` file. The controller is responsible for watching Foo resources and
reconciling the state of the system based on the observed state of the system.

The controller is defined as follows:

```go
// Controller is the controller implementation for Foo resources
type Controller struct {
// kubeclientset is a standard kubernetes clientset
kubeclientset kubernetes.Interface
// sampleclientset is a clientset for our own API group
sampleclientset clientset.Interface

deploymentsLister appslisters.DeploymentLister
deploymentsSynced cache.InformerSynced
foosLister        listers.FooLister
foosSynced        cache.InformerSynced

// workqueue is a rate limited work queue. This is used to queue work to be
// processed instead of performing it as soon as a change happens. This
// means we can ensure we only process a fixed amount of resources at a
// time, and makes it easy to ensure we are never processing the same item
// simultaneously in two different workers.
workqueue workqueue.TypedRateLimitingInterface[cache.ObjectName]
// recorder is an event recorder for recording Event resources to the
// Kubernetes API.
recorder record.EventRecorder
}

```

- kubeclientset is a standard kubernetes clientset
- sampleclientset is a clientset for our own API group, here is client of foo
- deploymentsLister is a lister for Deployment resources
- deploymentsSynced is a function that returns a boolean indicating whether the Deployment informer has synced
- foosLister is a lister for Foo resources
- foosSynced is a function that returns a boolean indicating whether the Foo informer has synced
- workqueue's explained as above
- recorder is an event recorder for recording Event resources to the Kubernetes API

controller real implement as follows:

```go
// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Foo resource
// with the current status of the resource.
func (c *Controller) syncHandler(ctx context.Context, objectRef cache.ObjectName) error {
logger := klog.LoggerWithValues(klog.FromContext(ctx), "objectRef", objectRef)

// Get the Foo resource with this namespace/name
foo, err := c.foosLister.Foos(objectRef.Namespace).Get(objectRef.Name)
if err != nil {
// The Foo resource may no longer exist, in which case we stop
// processing.
if errors.IsNotFound(err) {
utilruntime.HandleErrorWithContext(ctx, err, "Foo referenced by item in work queue no longer exists", "objectReference", objectRef)
return nil
}

return err
}

deploymentName := foo.Spec.DeploymentName
if deploymentName == "" {
// We choose to absorb the error here as the worker would requeue the
// resource otherwise. Instead, the next time the resource is updated
// the resource will be queued again.
utilruntime.HandleErrorWithContext(ctx, nil, "Deployment name missing from object reference", "objectReference", objectRef)
return nil
}

// Get the deployment with the name specified in Foo.spec
deployment, err := c.deploymentsLister.Deployments(foo.Namespace).Get(deploymentName)
// If the resource doesn't exist, we'll create it
if errors.IsNotFound(err) {
deployment, err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Create(context.TODO(), newDeployment(foo), metav1.CreateOptions{FieldManager: FieldManager})
}

// If an error occurs during Get/Create, we'll requeue the item so we can
// attempt processing again later. This could have been caused by a
// temporary network failure, or any other transient reason.
if err != nil {
return err
}

// If the Deployment is not controlled by this Foo resource, we should log
// a warning to the event recorder and return error msg.
if !metav1.IsControlledBy(deployment, foo) {
msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
c.recorder.Event(foo, corev1.EventTypeWarning, ErrResourceExists, msg)
return fmt.Errorf("%s", msg)
}

// If this number of the replicas on the Foo resource is specified, and the
// number does not equal the current desired replicas on the Deployment, we
// should update the Deployment resource.
if foo.Spec.Replicas != nil && *foo.Spec.Replicas != *deployment.Spec.Replicas {
logger.V(4).Info("Update deployment resource", "currentReplicas", *foo.Spec.Replicas, "desiredReplicas", *deployment.Spec.Replicas)
deployment, err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Update(context.TODO(), newDeployment(foo), metav1.UpdateOptions{FieldManager: FieldManager})
}

// If an error occurs during Update, we'll requeue the item so we can
// attempt processing again later. This could have been caused by a
// temporary network failure, or any other transient reason.
if err != nil {
return err
}

// Finally, we update the status block of the Foo resource to reflect the
// current state of the world
err = c.updateFooStatus(foo, deployment)
if err != nil {
return err
}

c.recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
return nil
}
```

as we can see in the code, we get the foo object from the cache, then get the deployment object from the cache, if the
deployment object does not exist, we will create it,
if the deployment object is not controlled by the foo object, we will return an error, if the number of replicas in the
foo object is different from the number of replicas in the deployment object,
we will update the deployment object, finally, we will update the status of the foo object.

## controller manager

in this repo, we did not implement a manager, but we can watch and sync other resources in other controller, we can
implement a manager to manage all the controllers in the run.go.
