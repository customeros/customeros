Feature: Contract Read

  Scenario: Custom SLIs are properly inserted
#    Given a tenant was created
#    And the organization Test Organization was created
#    And a contract was created for the organization Test Organization
    When the following SLIs are inserted in the database
      | billingType | price | quantity | startedAt                      |
      | MONTHLY     | 12    |2         | 2013-01-01T00:00:00.000000000Z |
      | ONCE        | 10    |3         | 2013-03-31T23:59:59.999999999Z |
    Then the SLIs should exist in the neo4j database in a consistent format

  Scenario Outline: Default SLIs are properly inserted
#    Given <number_of_sli> SLIs were inserted in the database
    When <inserted_slis> SLIs are inserted in the database
    Then <actual_number_of_SlIs> should exist in the neo4j database

    Examples:
      | inserted_slis | actual_number_of_SlIs |
      | 2            | 2         |
      | 1            | 1         |
