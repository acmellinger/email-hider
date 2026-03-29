export default {
  async fetch(request, env, ctx) {
    // 1. Define CORS Headers
    const corsHeaders = {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "POST, OPTIONS",
      "Access-Control-Allow-Headers": "Content-Type",
      "Access-Control-Max-Age": "3600",
    };

    // 2. Handle OPTIONS (Preflight)
    if (request.method === "OPTIONS") {
      return new Response(null, {
        status: 204,
        headers: corsHeaders,
      });
    }

    // 3. Log Origin
    const origin = request.headers.get("origin");
    if (origin) console.log(origin);

    // 4. Validate Method
    if (request.method !== "POST") {
      console.log(`Error: ${request.method} not allowed.`);
      return new Response("Only POST requests are accepted.", {
        status: 405,
        headers: corsHeaders,
      });
    }

    // 5. Parse JSON Body
    let body;
    try {
      body = await request.json();
    } catch (err) {
      console.log("Error: " + err.message);
      return new Response("Bad request body.", {
        status: 400,
        headers: corsHeaders,
      });
    }

    // 6. Validate Turnstile Token
    if (!body.Token) {
      console.log("Error: token not provided");
      return new Response("No Turnstile token provided.", {
        status: 400,
        headers: corsHeaders,
      });
    }

    // Cloudflare Turnstile Verify Endpoint
    const verifyUrl = "https://challenges.cloudflare.com/turnstile/v0/siteverify";
    
    const formData = new FormData();
    formData.append("secret", env.TURNSTILE_SECRET);
    formData.append("response", body.Token);
    // Passing the client's IP is good practice for Turnstile
    formData.append("remoteip", request.headers.get("CF-Connecting-IP"));

    let turnstileResult;
    try {
      const turnstileResp = await fetch(verifyUrl, {
        method: "POST",
        body: formData,
      });
      if (!turnstileResp.ok) {
        throw new Error(`Turnstile API returned ${turnstileResp.status}`);
      }
      turnstileResult = await turnstileResp.json();
    } catch (err) {
      console.log("Error: Turnstile API request failed: " + err.message);
      return new Response("Validation service error.", {
        status: 503,
        headers: corsHeaders,
      });
    }

    if (!turnstileResult.success) {
      console.log("Error: Turnstile validation failed", turnstileResult["error-codes"]);
      return new Response("Validation failed.", {
        status: 403,
        headers: corsHeaders,
      });
    }

    // 7. Validate Site and Retrieve Email
    if (body.Site) {
      const escapedSite = body.Site.replace(/[\n\r]/g, "");

      // Prefix lookup ensures only SITE_* env vars are accessible,
      // preventing injection of secret variable names (e.g. TURNSTILE_SECRET).
      const email = env["SITE_" + escapedSite];

      if (!email) {
        console.log(`Error: ${escapedSite} site not available`);
        return new Response("Validation failed.", {
          status: 403,
          headers: corsHeaders,
        });
      }

      console.log(`Email retrieved for site: ${escapedSite}`);

      return new Response(JSON.stringify({ email: email }), {
        status: 200,
        headers: {
          ...corsHeaders,
          "Content-Type": "application/json",
        },
      });

    } else {
      console.log("Error: site not provided");
      return new Response("No site provided.", {
        status: 400,
        headers: corsHeaders,
      });
    }
  },
};