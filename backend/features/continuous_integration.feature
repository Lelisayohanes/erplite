Feature: Continuous Integration
  As a development team member
  I want automated validation for every code change
  So that software quality is maintained before changes are merged or deployed

  Scenario: Validate changes on every push
    Given code is pushed to the repository
    When the CI pipeline executes
    Then code quality checks are performed
    And automated tests are executed
    And the pipeline completes successfully

  Scenario: Validate pull requests
    Given a pull request is opened
    When the CI pipeline executes
    Then the application is built successfully
    And the container image is validated

  Scenario: Improve pipeline performance
    Given dependencies have previously been downloaded
    When the pipeline executes
    Then dependency caching is used to reduce execution time
