// Types and interfaces
interface TokenValidator {
    validateToken(token: string): Promise<boolean>;
}
  
interface ValidationResponse {
    is_valid: boolean;
}
  
  // Token validator implementation
export class HankoTokenValidator implements TokenValidator {
    constructor(private readonly hankoApiUrl: string) {}
  
    async validateToken(token: string): Promise<boolean> {
      if (!token || token.length === 0) {
        return false;
      }
  
      try {
        const response = await fetch(`${this.hankoApiUrl}/sessions/validate`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ session_token: token }),
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
}