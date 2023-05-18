# define SEED_PIN A3

void setup() {
  Serial.begin(9600);
}

void loop() {
  int recdata = Serial.read();
  if (recdata != -1)
  {
    randomSeed(analogRead(SEED_PIN));
    int randNum = random(1000, 10000);
    Serial.println(randNum);
  }
}
