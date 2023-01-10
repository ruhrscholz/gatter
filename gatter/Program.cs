using System.Data;
using Npgsql;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.

// TODO Read connectionString from file
builder.Services.AddSingleton<IDbConnection>((sp) => new NpgsqlConnection("Host=localhost;Database=gatter"));
builder.Services.AddControllers();
builder.Services.AddRazorPages();

var app = builder.Build();

app.UseHttpsRedirection();
app.UseStaticFiles();

app.UseRouting();

app.UseAuthorization();

app.MapControllers();
app.MapRazorPages();

app.Run("http://localhost:8000");
