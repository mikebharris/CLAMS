@Regression
Feature: Attendee Writer service processes incoming records from BAMS and updates CLAMS

  @happy_path
  Scenario: New records from BAMS are added to CLAMS
    # New records are stored
    When the Attendee Writer is invoked with an attendee record from BAMS to be processed
    Then an attendee record is added to CLAMS

  Scenario: Records in CLAMS are updated with changes to records in BAMS
    When the Attendee Writer is invoked with an updated attendee record from BAMS to be processed
    Then the attendee record is updated in CLAMS
