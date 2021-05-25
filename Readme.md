# Credentials Rotator


## Use Case

Generally when provisioning infrastructure via something like terraform in
a GitOps Workflow your CI/CD lives outside of said infrastructure. This creates
a small issue, as for each CI/CD that lives outside of your system you would
need to have some kind of key rotation in place. This is a simple application
that handles this (currently specificly for Google Cloud).

### How it works

Currently it will create a new key under a Service Account update the repos variable and then delete all other keys.

## Example Config

```yaml
credentials:
- type: gitlab
  project_id: 12344
  variable: GOOGLE_CLOUD_CREDENTIALS
  service_account: example-1234@super-awesome-project.google.com
```
