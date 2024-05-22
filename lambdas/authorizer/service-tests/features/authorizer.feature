Feature: Checks username and password

  Scenario: Valid credentials return success payload
    When The lambda is invoked with valid credentials
    Then A Success response is returned

  Scenario: Invalid credentials return fail payload
    When The lambda is invoked with invalid credentials
    Then A Failure response is returned
