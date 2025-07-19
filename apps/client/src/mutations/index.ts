import { MutationOptions } from "@tanstack/react-query";

export type MutationOptionsHelper<TData, TVariables> = Omit<
  MutationOptions<TData, Error, TVariables>,
  "mutationKey" | "mutationFn"
>;

export const fetchPostRequest = async <T>(
  path: string,
  body: Record<string, any>
): Promise<T> => {
  const response = await fetch(
    `${import.meta.env.VITE_BACKEND_API_URL}${path}`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }
  );

  if (!response.ok) {
    throw new Error(
      `Failed to fetch ${path}: ${response.status} ${response.statusText}`
    );
  }

  const data = await response.json();
  return data as T;
};
