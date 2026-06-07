use actix_web::{App, HttpResponse, HttpServer, Responder, get, post, web::{self,  Query, Json}};
use serde::{ Serialize, Deserialize};

// #[derive(Deserialize)]
#[derive(Deserialize)]
struct  AddQuery{
    a: i32,
    b: i32,
}


#[derive(Debug, Serialize)]
struct AddResponse {
    result: i32,
}

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello world!")
}

#[post("/echo")]
async fn echo(req_body: String) -> impl Responder {
    HttpResponse::Ok().body(req_body)
}

// #[get("/steve")]
async  fn steve(query: Query<AddQuery>) -> Json<AddResponse>{
    let sum: i32 = query.a + query.b;
    // HttpResponse::Ok().body(sum.to_string())
    Json(AddResponse{result:sum})
    // takes in two numbers and adds them together
}

async fn manual_hello() -> impl Responder {
    HttpResponse::Ok().body("Hey there!")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(hello)
            .service(echo)
            // .service(steve)
            .route("/hey", web::get().to(manual_hello))
            .route("/steve", web::get().to(steve))
            // .route("/steve", web::post().to(steve))
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
