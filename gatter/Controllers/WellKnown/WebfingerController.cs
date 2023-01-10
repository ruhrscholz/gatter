using System.Data;
using Microsoft.AspNetCore.Mvc;
using Dapper;
using gatter.Model;

namespace gatter.Controllers.WellKnown;

[ApiController]
[Route("/.well-known/webfinger")]
public class WebfingerController : ControllerBase
{
    [HttpGet]
    [ProducesResponseType(StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public ActionResult<User> Get(
	    IConfiguration configuration,
	    IDbConnection dbConnection,
	    [FromQuery]string resource)
    {
	    
	    resource = resource.StartsWith("acct:") ?
		    resource.Remove(0, 5) :
		    resource;

	    string[] resourceSplit = resource.Split("@");

	    if (resourceSplit.Length > 2)
	    {
		    return StatusCode(StatusCodes.Status400BadRequest, string.Empty);
	    }
	    
	    if (resourceSplit.Length == 2 &&
	        !resourceSplit[1].Equals(configuration.GetSection("Domains").GetValue<string>("Web"), StringComparison.OrdinalIgnoreCase) &&
	        !resourceSplit[1].Equals(configuration.GetSection("Domains").GetValue<string>("Local"), StringComparison.OrdinalIgnoreCase))
	    {
		    return StatusCode(StatusCodes.Status404NotFound, string.Empty);
	    }

	    User user = dbConnection.QueryFirstOrDefault<User>("select Username=@username", new { username = resourceSplit[0] });

	    return user;
	    
	    throw new NotImplementedException();
        
        /*
         *
         
		// TODO Exists
		var username string
		if err := rows.Scan(&username); err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			log.Printf("Could not query database for webfinger: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/jrd+json")

		response := webfingerResponse{
			Subject: fmt.Sprintf("acct:%s@%s", username, env.LocalDomain),
			Aliases: []string{
				fmt.Sprintf("https://%s/@%s", env.WebDomain, username),
				fmt.Sprintf("https://%s/users/%s", env.WebDomain, username),
			},
			Links: []webfingerResponseLink{
				{
					Rel:   "http://webfinger.net/rel/profile-page",
					Type_: "text/html",
					Href:  fmt.Sprintf("https://%s/@%s", env.WebDomain, username),
				},
				{
					Rel:   "self",
					Type_: "application/activity+json",
					Href:  fmt.Sprintf("https://%s/users/%s", env.WebDomain, username),
				},
				{
					Rel:      "http://ostatus.org/schema/1.0/subscribe",
					Template: fmt.Sprintf("https://%s/authorize_interaction?uri={uri}", env.WebDomain),
				},
			},
		}

		json.NewEncoder(w).Encode(response)
         */
    }
    
    
    public class WebfingerResponse
    {
        private string subject;
        private string[] aliases;
        private Link[] links;

        class Link
        {
            
        }
    }
}
