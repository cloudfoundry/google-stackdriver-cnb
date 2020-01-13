# `google-stackdriver-cnb`
The Cloud Foundry Google Stackdriver Buildpack is a Cloud Native Buildpack V3 that provides the [Google Stackdriver][g] Debugger and Profiler agents and configuration to applications.

This buildpack is designed to work in collaboration with bound service instances.

[g]: https://cloud.google.com/stackdriver/

## Behavior
This buildpack will participate if either of the following conditions are met

* A service is bound with a payload containing `binding_name`, `instance_name`, `label`, or `tag` containing `google-stackdriver-debugger` as a substring.
* A service is bound with a payload containing `binding_name`, `instance_name`, `label`, or `tag` containing `google-stackdriver-profiler` as a substring.

The buildpack will do the following if the debugger is bound:

* Contributes Google Stackdriver Debugger agent to a layer marked `launch`
* Sets `-agentpath` and `com.google.cdbg.auth.serviceaccount.enable` to `$JAVA_OPTS`
* Contributes Google Stackdriver Credentials helper to a layer marked `launch`.
* Sets `$GOOGLE_APPLICATION_CREDENTIALS`.
  
The buildpack will do the following if the profiler is bound:
* Contributes Google Stackdriver Profiler agent to a layer marked `launch`
* Sets `-agentpath` to `$JAVA_OPTS`
* Contributes Google Stackdriver Credentials helper to a layer marked `launch`.
* Sets `$GOOGLE_APPLICATION_CREDENTIALS`.

## Configuration 
| Environment Variable | Description
| -------------------- | -----------
| `$BPL_GOOGLE_STACKDRIVER_MODULE` | Configures the module name. Defaults to `default-module`.
| `$BPL_GOOGLE_STACKDRIVER_VERSION` | Configures the module version. Defaults to empty.

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: https://www.apache.org/licenses/LICENSE-2.0

