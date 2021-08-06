# Sendgrid API

This script is needed to add new users to SendGrid as SSO teammates.
Previously, all users were manually added and manually migrating them would take a long time.

### Preparation for run script
##### Users config
Update users in `conf.yaml` file, add the user to one of three groups, they all differ in access rights.
```yaml
groups:
  - admin:
    users:
      - email: admin@example.com (
        first_name: Zack
        last_name: Yo
    scopes: []
  - developer:
    users:
      - email: developer@example.com
        first_name: Nik
        last_name: Nilson
    scopes: []
  - support:
    users:
      - email: support@example.com
        first_name: Vladimir
        last_name: Nikolaev
    scopes: []
```

##### Build
```bash
env GO111MODULE=on go build -o sg-users-sso -v
```

##### Run
```bash
Create all users:
env SENDGRID_API_KEY="xxx" ./sg-users-sso --create

Update all users:
env SENDGRID_API_KEY="xxx" ./sg-users-sso --update

Get all no SSO users:
env SENDGRID_API_KEY="xxx" ./sg-users-sso --get-all-no-sso

Get all users:
env SENDGRID_API_KEY="xxx" ./sg-users-sso --get-all

See help:
./sg-users-sso --help
```
> Also you can use `--sendgrid-token` for setting SENDGRID_API_KEY
> ```bash
> ./sg-users-sso --get-all --sendgrid-token="xxx"
> ```

##### Available flags:
```bash
      --config-path string      Config file path (default "config/users.yaml")
  -c, --create                  Create all users
  -a, --get-all                 Get all users
  -n, --get-all-no-sso          Get all no SSO users
  -t, --sendgrid-token string   Config file path, default env: SENDGRID_API_KEY
  -u, --update                  Update all users
```

### Access scopes:
Due to the fact that SendGrid does not support groups, you have to manage separate scopes.
You must grant the necessary permissions to each group, all available scopes are in the file `config/users.yaml`
1. admin group (access to all scopes);
2. developers group;
3. support group.


# Licence
Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements. See the NOTICE file distributed with this work for additional information regarding copyright ownership. The ASF licenses this file to you under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
