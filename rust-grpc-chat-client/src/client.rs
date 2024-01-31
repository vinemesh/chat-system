use std::sync::mpsc;
use std::io::{self, Write};
use tonic::transport::Endpoint;
use crate::message_service::message_service_client::MessageServiceClient;
use crate::message_service::PlayerMessage;
use tower::util::ServiceExt;

// Incluindo os códigos gerados a partir do arquivo .proto
pub mod message_service {
    tonic::include_proto!("messages");
}

// Função assíncrona para enviar mensagens através do stream
async fn listen_for_messages(mut client: MessageServiceClient<tonic::transport::Channel>, rx: mpsc::Receiver<String>) -> Result<(), Box<dyn std::error::Error>> {
    let mut inbound = client.stream_messages(tonic::Request::new(async_stream::stream! {
        while let Ok(msg) = rx.recv() {
            // Enviando cada mensagem recebida pelo canal como parte do stream
            yield PlayerMessage { player_id: "123".to_string(), content: msg };
        }
    })).await?.into_inner();

    // Processando mensagens recebidas do servidor
    while let Some(message) = inbound.message().await? {
        println!("Received: {:?}", message.message);
    }

    Ok(())
}

#[tokio::main]
async fn main() {
    // Parse dos argumentos da linha de comando
    let args: Vec<String> = std::env::args().collect();
    let endpoint: String = args.get(1).unwrap_or(&"http://localhost:50051".to_string()).clone();

    let endpoint = endpoint.parse::<Endpoint>().expect("Invalid endpoint");
    let mut channel = endpoint.connect_lazy();

    // Verificando a conexão com o servidor
    if let Err(e) = channel.ready().await {
        eprintln!("Falha ao conectar ao servidor: {}", e);
        return;
    }
    println!("Conectado ao servidor gRPC.");

    let client = MessageServiceClient::new(channel.clone());

    // Criando um canal para comunicação entre o loop de entrada e a tarefa de streaming
    let (tx, rx) = mpsc::channel();

    // Iniciando a tarefa de streaming em uma nova task assíncrona
    let client_clone = client.clone();
    tokio::spawn(async move {
        listen_for_messages(client_clone, rx).await.expect("Failed to listen for messages");
    });

    // Loop para ler mensagens do terminal e enviá-las para o servidor
    loop {
        print!("Enter message: ");
        io::stdout().flush().unwrap(); // Garante que "Enter message:" seja exibido imediatamente

        let mut content = String::new();
        io::stdin().read_line(&mut content).expect("Failed to read line");
        let content = content.trim().to_string(); // Removendo espaços em branco e nova linha

        if content == "/quit" {
            break;
        }

        // Enviar a mensagem para o canal, que será processado pela tarefa de streaming
        tx.send(content).expect("Failed to send message to the stream");
    }
}
