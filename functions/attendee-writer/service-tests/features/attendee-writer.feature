@Regression
Feature: Attendee Writer service processes results and updates Attendees datastore

  @happy_path
  Scenario: The attendees datastore is kept up to date with the records in BAMS
    # New records are stored
    When the Attendee Writer is invoked with an attendee record from BAMS to be processed
    Then the attendee is added to the Attendees Datastore

    # Existing records are updated
    When the Attendee Writer is invoked with an updated attendee record from BAMS to be processed
    Then the attendee is updated in the Attendees Datastore
