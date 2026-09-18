import type { UserResponse, User } from '@/lib/api/types';
import { serviceApi } from '@/lib/api/endpoints';

export default async function MePage() {
  const data: UserResponse = await serviceApi.me();
  // TODO: catch api error (maybe in error.tsx component)
  const user: User | undefined = data.user;

  if (!user) {
    return <div>User not found</div>;
  }

  return (
    <>
      <p>Email: {user.email}</p>
      <p>Role: {user.role}</p>
    </>
  );
}
