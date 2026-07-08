Feature: Dependency Injection
  As a developer
  I want application dependencies to be managed through dependency injection
  So that the application remains modular, maintainable, and easily testable

  Scenario: Resolve application dependencies
    Given application dependencies are registered
    When the application starts
    Then all required components are resolved successfully

  Scenario: Build the application using dependency injection
    Given the dependency container is configured
    When the application is initialized
    Then all application services are injected automatically
    And the application starts successfully
