using Microsoft.AspNetCore.Mvc;

namespace gatter.Controllers.WellKnown;

[ApiController]
[Route("/.well-known/webfinger")]
public class WebfingerController : ControllerBase
{
    [HttpGet]
    [ProducesResponseType(StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public ActionResult<WebfingerResponse> Get()
    {
        throw new NotImplementedException();
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
