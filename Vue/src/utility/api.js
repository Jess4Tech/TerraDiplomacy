export const API_URL = 'https://127.0.0.1:10000/api/v1';

/**
 * Get user authorization data
 * @return {{bool, number}}
 */
export async function getAuthorized() {
  const res = await fetch(`${API_URL}/auth/status`, {
    method: 'GET',
    mode: 'cors',
    credentials: 'include',
  });
  try {
    return (await res.json());
  } catch (e) {
    console.log(e);
    return {auth: false, tier: 0};
  }
}

/**
 * Executes the success or fail function if authorized for the specified tier
 * (Specified functions default to nothing)
 * @param {function} success - Function to execute if authorized
 * @param {function} fail - Function to execute if not authorized
 * @param {number} tier - Needed authorization tier
 */
export async function ifAuthorized(
    success = () => {},
    fail = () => {},
    tier = 1,
) {
  const authorized = await getAuthorized();
  if (authorized.auth && authorized.tier >= tier) {
    success(authorized);
  } else {
    fail(authorized);
  }
}

/**
 * Redirect to login if the user is unauthorized
 * @param {VueRouter} router - Vue Router
 * @param {number} tier - Needed authorization tier
 */
export async function redirectIfUnauthorized(router, tier = 1) {
  ifAuthorized(() => {}, async () => {
    console.log('Not logged in, redirecting');
    await router.push({name: 'LoginPage'});
  }, tier);
}

/**
 * Delete project
 * @param {string} name - The name of the project to remove
 * @return {bool}
 */
export async function deleteProject(name) {
  return (await fetch(`${API_URL}/projects`, {
    method: 'DELETE',
    mode: 'cors',
    credentials: 'include',
    body: JSON.stringify({
      name: name,
    }),
  })).status == 200;
}

/**
 * Get all projects
 * @return {[{string, string, number}]}
 */
export async function getProjects() {
  const res = await fetch(`${API_URL}/projects`, {
    method: 'GET',
    mode: 'cors',
    credentials: 'include',
  });
  try {
    return await res.json();
  } catch (e) {
    console.log(`fetching failed`, e);
    return [];
  }
}

/**
 * Get all tension entries
 * @return {[number, number]}
 */
export async function getTension() {
  const res = await fetch(`${API_URL}/tension`, {
    method: 'GET',
    mode: 'cors',
    credentials: 'include',
  });
  try {
    return await res.json();
  } catch (e) {
    console.log(`fetching failed`, e);
    return [];
  }
}
