Feature: Workshop Signups Processor processes workshop signups

  Scenario: Workshop Signups Processor logs signups in CLAMS and notifies workshop leaders
    Given workshop signup records exist in the database
    When the workshop signup processor receives a notification
    Then the CLAMS datastores are updated
    And a report for the workshop facilitator is produced
    And and the workshop facilitator is notified by email