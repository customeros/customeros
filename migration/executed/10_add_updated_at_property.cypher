# execute after releasing https://github.com/openline-ai/openline-customer-os/issues/649

match (n) where n.updatedAt is null AND n.createdAt is not null with n SET n.updatedAt=n.createdAt;