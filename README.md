# Pictura Certamine

## Features

- [x] form with user data nad multiple pictures
- [x] confirmation email
- [x] saving progress of form
- [x] sending errors to sentry
- [x] manage cookies (only required, used for performance and all)
- [x] handle duplicate email
- [x] generic 404 page
- [x] list view crm
- [x] files download
- [ ] export to excel/csv
- [ ] admin user that can view crm
- [ ] translations
- [ ] typography
- [ ] UI form looks
- [ ] captcha
- [ ] what reverse proxy
- [ ] backup sqlite to s3 (research how to)
- [ ] (optional) maybe a rating module
- [ ] (optional) add logging
- [ ] (optional) add metrics

## External Systems

- [x] s3 free tier at oracle cloud
- [x] email (sendgrid or other service which is cheaper)
  - sendgrid 100 emails/day (won because easy integration)
  - mailgun 100 emails/day (only 1 api key)
  - emaillabs 800 email/day (dunno if integration is hard)
  - mailhog for testing
- [ ] vm (vm on oracle should be fine)
- [ ] captcha (google reCaptcha v2 - I'm not a robot - easier to integrate then v3)
- [x] sentry (hope for a generous free trial - 5k errors - should be plenty)
- [ ] prometheus for metrics (free tier on grafana labs)
- [x] sqlite

## Tech Stack

- go
- postgres/sqlite
- htmx?
- gotempl
