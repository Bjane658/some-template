
@Slf4j
@Service
@AllArgsConstructor
@KafkaListener(
    topics = "account.public.inconsistency",
    autoStartup = "${sipgate.kafka.consumer.auto-startup}")
public class KafkaConsumer {

  private final TestService testService;

  @KafkaHandler
  public void consumerMessage(
      @Payload final Object someMessage) {
    MDC.put("masterSipId", someMessage.getMasterSipId());
    MDC.put("domain", someMessage.getDomain());

    log.info("Received message");

    MDC.clear();
  }

  @KafkaHandler(isDefault = true)
  public void handleDefault(@Payload ConsumerRecord<?, ?> record) {
    log.info("Ignoring unknown event.");
  }
}

