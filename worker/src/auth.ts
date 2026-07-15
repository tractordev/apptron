  
interface ValidationResponse {
    is_valid: boolean;
}

export async function validateToken(hankoApiUrl: string, token: string): Promise<boolean> {
  if (!token || token.length === 0) {
    return false;
  }

  try {
    const response = await fetch(`${hankoApiUrl}/sessions/validate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ session_token: token }),
      redirect: 'manual',
    });

    if (!response.ok) {
      return false;
    }

    const validationData = await response.json() as ValidationResponse;
    return validationData.is_valid;
  } catch (error) {
    console.error('Token validation error:', error);
    return false;
  }
}

export function parseJWT(token: string): Record<string, any> {
  const base64Url = token.split('.')[1];
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  return JSON.parse(atob(base64));
}