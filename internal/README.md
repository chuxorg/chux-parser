# internal
Within internal, packages are structured by features in order to achieve the so-called screaming architecture. 
For example, an `auth` directory contains the application logic related to the `auth` feature.

The packages are:
- `auth/`     - For code that has to do with authentication
- `config/`   - For code that has to do with configuration
- `groups/`   - For code that has to do with retrieving and formatting Okta Groups
- `interfaces/` - For Go interfaces for idc-okta-api
- `jks/`        - For code that has to do with signing the Bearer Token for accessing the Okta API
- `models/`     - For code that has to do with modeling HTTP responses
- `routes/`     - For code that has to do with http routing 
- `version/`    - For code that has to do with idc-okta-api versioning

