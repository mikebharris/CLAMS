Feature: Workshop Signups Processor processes workshop signups

  Scenario: Workshop Signups Processor logs signups in CLAMS and notifies workshop leaders
    Given workshop signup records exist in the database
    When the workshop signup processor receives a notification
    Then the workshops signups datastore is updated
