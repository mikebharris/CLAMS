Feature: CLAMS Attendees API Returns details of attendees

  Scenario: Attendees API returns list of all attendees to event
    Given some attendee records exist in the attendees datastore
    When the front-end requests a specific attendee record from the API
    Then the record is returned

    When the front-end requests all records from the API
    Then all available records are returned