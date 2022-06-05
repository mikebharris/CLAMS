Feature: CLAMS Report API Returns stats about the event

  Scenario: Report API Returns stats about the event
    Given some attendee records exist in the attendees datastore
    When the front-end requests the stats from the report API
    Then some statistics about the event are returned