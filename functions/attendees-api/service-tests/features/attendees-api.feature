Feature: CLAMS Attendees API Returns details of attendees

  Scenario: Attendees API returns list of all attendees to event
    Given an attendee record exists in the attendees datastore
    When the front-end requests the attendee record from the API
    Then the record is returned

    When the front-end requests all records from the API
    Then all available records are returned