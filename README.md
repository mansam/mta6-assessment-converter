# mta6-assessment-converter

## Usage:

`make cmd` to build command (requires go 1.21)

```sh
bin/convert-assessments path-to-assessments.json > output.json
```

Converts a file containing a JSON array of MTA6/Pathfinder assessment objects to an array of MTA7 format assessment objects.
The tool expects that the Pathfinder questionnaire is seeded at ID 1, which should be the case under normal circumstances.
Completion status, confidence, and risk assessment values are not transfered as these are dynamically determined by MTA7.

The output created by this tool requires a version of the Tackle import tool that is capable of importing MTA7 assessments.
If that is not possible, the converted assessments may be applied by POSTing each of them to the `/hub/applications/:id/assessments` endpoint
of a deployed Konveyor instance, where `:id` would be replaced by the assessment's application ID.

For example, an assessment that started like this:

```json
  {
    "createUser": "",
    "updateUser": "",
    "createTime": "0001-01-01T00:00:00Z",
    "application": {
      "id": 82
    },
    "questionnaire": {
      "id": 1
    },
    "sections": [
      {
        "order": 1,
        "name": "Application details",
        "questions": [
          {
            "order": 1,
<snip>
```
would be POSTed to `/hub/applications/82/assessments`.
