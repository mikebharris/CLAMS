Feature: Database Trigger

  Scenario: Database Trigger generates messages for each changed record
    Given there are database trigger records in the database
    When the database trigger is invoked
    Then the messages are placed on the queue
