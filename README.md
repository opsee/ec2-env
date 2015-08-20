# ec2-env

Generate exports for environment variables for AWS EC2 instances.

## Usage

Primarily useful in shell startup scripts and systemd units as an EnvironmentFile.

The variables exported are:

```
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
AWS_SESSION_TOKEN
AWS_INSTANCE_ID
AWS_IMAGE_ID
AWS_ACCOUNT_ID
AWS_DEFAULT_REGION
```

## Building

`$ gb build`
