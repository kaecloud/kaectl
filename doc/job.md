# config
kaectl needs a config file named `~/.kae/config.yaml`, below is an example content of config file

    sso_username: user
    sso_password: passwd
    sso_host: keyclock host
    sso_realm: kae
    sso_client_id: kae-cli

    kae_url: https://console.gtapp.xyz
    job_server_url: http://127.0.0.1:8080
    job_default_cluster: default cluster

the meaning of each field is clear.

# job spec

    name: jobname
    
    # +optional
    auto_restart: false
    
    # the fllowing fields are copy from k8s's JobSpec,
    # they have the same meaning as k8s.
    # see: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#jobspec-v1-batch
    # +optional
    parallelism: 3
    completions: 3
    activeDeadlineSeconds: 10  
    backoffLimit: 6
    ttlSecondsAfterFinished:
    
    # prepare is used to do some prepare operations, it mainly doees 2 things:
    # 1. download artifiact
    # 2. run prepare command
    prepare:
      artifacts:
          # url support mulitple protocols:
          # 1. http/https
          # 2. oss
          # 3. git 
        - url: xxx
          # For client: local is the path needs to upload to OSS
          # For server: local is the path we need download arfifact to
          local: dir1/dir2
          
      # the image used to run prepare command
      image: xxx
      # the prepare command
      command: xxx
      # if shell set to true, then command is run as a shell command
      # (sh -c command)
      shell: false
      
    # only specified cron when you want to create CronJob
    cron:
      # the fllowing fields are copy from k8s's CronJobSpec,
      # they have the same meaning as k8s.
      # see: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#cronjobspec-v1beta1-batch
      # required
      schedule: xxx
      startingDeadlineSeconds: xx
      concurrencyPolicy: xxx
      suspend: false
      successfulJobsHistoryLimit: xxx
      failedJobsHistoryLimit: xxx