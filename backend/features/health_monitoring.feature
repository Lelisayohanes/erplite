Feature: Health and Readiness Monitoring
  As a platform administrator
  I want health and readiness endpoints
  So that orchestration platforms can monitor application availability and readiness

  Scenario: Verify application liveness
    Given the application is running
    When the liveness endpoint is requested
    Then the application responds with HTTP 200
    And indicates it is healthy

  Scenario: Verify application readiness
    Given all required dependencies are available
    When the readiness endpoint is requested
    Then the application responds with HTTP 200
    And indicates it is ready to accept traffic

  Scenario: Report dependency failures
    Given a required dependency is unavailable
    When the readiness endpoint is requested
    Then the application responds with HTTP 503
    And the response indicates the application is not ready
