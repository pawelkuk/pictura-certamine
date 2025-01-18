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
- [x] export to excel/csv
- [x] admin user that can view crm
- [x] command that creates a link that enables to register a user
- [x] typography
- [ ] UI form looks (wip - add form in modal e.g. https://developer.mozilla.org/en-US/docs/Web/CSS/:modal)
  - [x] modal
  - [ ] add links to pdfs
  - [ ] change font to one which supports unicode characters
  - [x] footer
  - [x] show error message on contest form (or figure out why alert collides with dialog window)
  - [x] banner (when graphics are going to be provided)
  - [x] blur on cookie banner
- [ ] translations (i18n https://github.com/gin-contrib/i18n/blob/master/_example/main.go)
- [x] what reverse proxy (treafik)
- [x] setting up certificates
- [ ] backup sqlite to s3 (research how to; cronjob maybe)
- [x] setting up customer's email as sender
- [ ] deployment (almost ready)
- [ ] captcha
- [ ] (optional) maybe a rating module
- [ ] (optional) add logging
- [ ] (optional) add metrics
- [ ] (optional) tests

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
