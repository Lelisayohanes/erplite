Feature: Development Automation
  As a developer
  I want common development tasks to be automated through a Makefile
  So that development activities are standardized and efficient

  Scenario: Build the application
    Given the project source code exists
    When the build command is executed
    Then the application binary is generated successfully

  Scenario: Execute automated tests
    Given automated tests exist
    When the test command is executed
    Then all tests are run
    And the results are displayed

  Scenario: Validate code quality
    Given the linting tool is configured
    When the lint command is executed
    Then coding standard violations are reported

  Scenario: Execute database migrations
    Given the database is available
    When the migration command is executed
    Then all pending migrations are applied successfully

  Scenario: Display available development commands
    When the help command is executed
    Then all supported commands are displayed with descriptions
