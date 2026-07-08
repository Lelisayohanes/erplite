Feature: Structured Logging
  As a developer
  I want application logs to be structured and traceable
  So that I can efficiently monitor, troubleshoot, and diagnose production issues

  Scenario: Generate structured log entries
    Given the application is running
    When an event is logged
    Then the log entry is generated in a structured format
    And includes the timestamp, severity level, and message

  Scenario: Include correlation identifiers
    Given a request includes a correlation identifier
    When the request is processed
    Then all related log entries include the same correlation identifier

  Scenario: Log incoming HTTP requests
    Given request logging is enabled
    When an HTTP request is processed
    Then the request metadata is logged
    And sensitive information is excluded
