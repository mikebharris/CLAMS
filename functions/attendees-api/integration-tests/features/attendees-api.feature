Feature: CLAMS Attendees API Returns details about attendees

  Scenario: Attendees endpoint returns list of all attendees to event
    Given some attendee records exist in CLAMS
    When the front-end requests all records from the endpoint
    Then all available records are returned

  Scenario: A specific attendee can be requested from CLAMS
    Given some attendee records exist in CLAMS
    When the front-end requests a specific attendee record from the endpoint
    Then the record is returned

  Scenario: Report endpoint Returns stats about the event
    Given some attendee records exist in CLAMS
    When the front-end requests the stats from the report endpoint
    Then some statistics about the event are returned