use kafka::consumer::{Consumer, FetchOffset, GroupOffsetStorage};
use kafka::producer::{Producer, Record};
use std::env;
use tokio::sync::mpsc;
use tokio::task;
use dotenv::dotenv;
use std::sync::{Arc, Mutex};

async fn consume_and_forward(
    brokers: String,
    group: String,
    topic: String,
    tx: mpsc::Sender<Vec<u8>>,
) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let consumer = Arc::new(Mutex::new(
        Consumer::from_hosts(vec![brokers.clone()])
            .with_topic(topic.clone())
            .with_group(group.clone())
            .with_fallback_offset(FetchOffset::Earliest)
            .with_offset_storage(GroupOffsetStorage::Kafka)
            .create()
            .unwrap(),
    ));

    tokio::task::spawn_blocking(move || {
        let consumer = Arc::clone(&consumer);
        loop {
            let mss = consumer.lock().unwrap().poll().unwrap();
            for ms in mss.iter() {
                for m in ms.messages() {
                    if let Err(_) = tx.blocking_send(m.value.to_vec()) {
                        eprintln!("Receiver dropped");
                    }
                }
                consumer.lock().unwrap().consume_messageset(ms).unwrap();
            }
            consumer.lock().unwrap().commit_consumed().unwrap();
        }
    }).await?;

    Ok(())
}

async fn run_producer(brokers: String, topic: String, mut rx: mpsc::Receiver<Vec<u8>>) {
    let mut producer = Producer::from_hosts(vec![brokers])
        .create()
        .unwrap();

    while let Some(message) = rx.recv().await {
        let record = Record::from_value(&topic, message);
        producer.send(&record).unwrap();
    }
}


#[tokio::main]
async fn main() {
    dotenv().ok();
    
    let consume_brokers = env::var("CONSUME_BROKERS").expect("CONSUME_BROKERS must be set");
    let consume_topic = env::var("CONSUME_TOPIC").expect("CONSUME_TOPIC must be set");
    let procude_brokers = env::var("PRODUCE_BROKERS").expect("PRODUCE_BROKERS must be set");
    let produce_topic = env::var("PRODUCE_TOPIC").expect("PRODUCE_TOPIC must be set");
    let env = env::var("ENV").expect("ENV must be set");
    
    let group = format!("{}-group-{}", consume_topic, env);

    let (tx, rx) = mpsc::channel(1024);

    let consumer_task = task::spawn(consume_and_forward(consume_brokers, group, consume_topic, tx));
    let producer_task = task::spawn(run_producer(procude_brokers, produce_topic, rx));

    let _ = tokio::try_join!(consumer_task, producer_task);
}