Feature: Contract Read

  Scenario Outline: SLIs are properly inserted
#    Given <number_of_sli> SLIs were inserted in the database
    When <inserted_slis> SLIs were inserted in the database
    Then <actual_number_of_SlIs> should exist in the neo4j database

    Examples:
      | inserted_slis | actual_number_of_SlIs |
      | 2            | 2         |
      | 1            | 1         |
      | 3            | 2         |
