Feature: Example

  Scenario Outline: eat godogs
    Given there are <initial_count> godogs
    When I eat <eat_count>
    Then there should be <remaining_count> remaining

    Examples:
      | initial_count | eat_count | remaining_count |
      | 12            | 5         | 7               |
      | 13            | 6         | 7               |
      | 14            | 1         | 4               |