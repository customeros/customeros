import React from 'react';
import {useCreateOrganization, useOrganizationDetails} from "../../../../hooks/useOrganization";
import styles from "../organization-details.module.scss";
import {Link} from "../../../ui-kit";

interface OrganizationFormProps {

}

export const OrganizationCreate: React.FC<OrganizationFormProps> = () => {
    const { onCreateOrganization } = useCreateOrganization( );

    return (
        <div className={styles.organizationDetails}>
            <div className={styles.bg}>
                <div>
                    <h1 className={styles.name}>{data?.name}</h1>
                    <span className={styles.industry}>{data?.industry}</span>
                </div>

                <p className={styles.description}>{data?.description}</p>

                {data?.website && <Link href={data.website}> {data.website} </Link>}
            </div>
        </div>
};
