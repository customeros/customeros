Feature: Contract Read

  @tag_custom_sLIs_are_properly_inserted
  Scenario: Custom SLIs are properly inserted
    Given a tenant was created
    And the organization with the id d1752d4b-cf87-474b-9f6c-b2736f26d977 was created
    And a contract with the id e1b9eac5-4d11-46dd-bb94-c9c7aa876f6c was created for the organization having the id d1752d4b-cf87-474b-9f6c-b2736f26d977
    When the following SLIs are inserted in the database
      | billingType | price | quantity | startedAt                      | contractId                           |
      | MONTHLY     | 12    |2         | 2013-01-01T00:00:00.000000000Z | e1b9eac5-4d11-46dd-bb94-c9c7aa876f6c |
      | ONCE        | 10    |3         | 2013-03-31T23:59:59.999999999Z | e1b9eac5-4d11-46dd-bb94-c9c7aa876f6c |
    Then the SLIs should exist in the neo4j database in a consistent format

  @tag_default_sLIs_are_properly_inserted
  Scenario Outline: Default SLIs are properly inserted
    Given a tenant was created
    And the organization with the id a55e0812-ba4b-4f3d-853b-73323cb3cdd6 was created
    And a contract with the id f181e1fb-3675-427c-9d4b-ebd61386d4ad was created for the organization having the id a55e0812-ba4b-4f3d-853b-73323cb3cdd6
    When <inserted_slis> SLIs are inserted in the database for the contract f181e1fb-3675-427c-9d4b-ebd61386d4ad
    Then <actual_number_of_SlIs> should exist in the neo4j database

    Examples:
      | inserted_slis | actual_number_of_SlIs |
      | 2            | 2         |
      | 1            | 13         |
