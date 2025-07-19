import { fetchPostRequest, MutationOptionsHelper } from ".";

import { useMutation } from "@tanstack/react-query";

export interface DatabaseCredentials {
  host: string;
  port: number;
  username: string;
  password: string;
  database: string;
}

export type NewConnectionOutput = {
  status: "connected";
  databaseId: string;
};

export const useNewConnectionMutation = (
  options?: MutationOptionsHelper<NewConnectionOutput, DatabaseCredentials>
) => {
  const path = "/new-connection";

  return useMutation({
    ...options,
    mutationKey: [path],
    mutationFn: async (credentials) => {
      const response = await fetchPostRequest<NewConnectionOutput>(
        path,
        credentials
      );
      return response;
    },
  });
};
