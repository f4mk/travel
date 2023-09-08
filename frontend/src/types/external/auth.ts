/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */

export interface paths {
  '/auth/login': {
    post: {
      requestBody: {
        content: {
          'application/json': components['schemas']['LoginUser']
        }
      }
      responses: {
        /** @description user logged in successfully */
        201: {
          content: {
            'application/json': components['schemas']['UserResponse']
          }
        }
        /** @description bad request */
        400: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description unauthorized */
        401: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description internal server error */
        500: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
      }
    }
  }
  '/auth/logout': {
    post: {
      requestBody: {
        content: {
          'application/json': Record<string, never>
        }
      }
      responses: {
        /** @description user logged in successfully */
        201: {
          content: {
            'application/json': Record<string, never>
          }
        }
        /** @description bad request */
        400: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description unauthorized */
        401: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description not found */
        404: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description internal server error */
        500: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
      }
    }
  }
  '/auth/password/change': {
    post: {
      requestBody: {
        content: {
          'application/json': components['schemas']['ChangePassword']
        }
      }
      responses: {
        /** @description password is changed successfully */
        201: {
          content: {
            'application/json': components['schemas']['UserResponse']
          }
        }
        /** @description bad request */
        400: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description unauthorized */
        401: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description not found */
        404: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description internal server error */
        500: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
      }
    }
  }
  '/auth/password/reset': {
    post: {
      requestBody: {
        content: {
          'application/json': components['schemas']['ResetPassword']
        }
      }
      responses: {
        /** @description password reset request was handled successfully */
        201: {
          content: {
            'application/json': Record<string, never>
          }
        }
        /** @description bad request */
        400: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description internal server error */
        500: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
      }
    }
  }
  '/auth/password/reset/submit': {
    post: {
      requestBody: {
        content: {
          'application/json': components['schemas']['SubmitResetPassword']
        }
      }
      responses: {
        /** @description password is changed successfully */
        201: {
          content: {
            'application/json': components['schemas']['UserResponse']
          }
        }
        /** @description bad request */
        400: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description bad request */
        403: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
        /** @description internal server error */
        500: {
          content: {
            'application/json': components['schemas']['ErrorResponse']
          }
        }
      }
    }
  }
}

export type webhooks = Record<string, never>

export interface components {
  schemas: {
    LoginUser: {
      /** @description user email */
      email: string
      /** @description user password */
      password: string
    }
    UserResponse: {
      /** @description user unique id */
      id: string
      /** @description user's name */
      name: string
      /** @description user's email */
      email: string
      /** @description date created */
      date_created: string
    }
    ResetPassword: {
      /** @description user email */
      email: string
    }
    ErrorResponse: {
      /** @description error message */
      error: string
      fields?: {
        [key: string]: string
      }
    }
    SubmitResetPassword: components['schemas']['ResetToken'] &
      components['schemas']['NewPassword']
    ChangePassword: components['schemas']['OldPassword'] &
      components['schemas']['NewPassword']
    OldPassword: {
      /** @description user new password */
      password_old: string
    }
    NewPassword: {
      /** @description user new password */
      password: string
      /** @description user new password confirm */
      password_confirm: string
    }
    ResetToken: {
      /** @description password reset secret token */
      token: string
    }
  }
  responses: never
  parameters: never
  requestBodies: never
  headers: never
  pathItems: never
}

export type external = Record<string, never>

export type operations = Record<string, never>
