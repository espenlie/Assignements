#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t lock;

// Note the return type: void*
void* adder(){
    for(int x = 0; x < 1000000; x++){
        pthread_mutex_lock(&lock);
        i++;
        pthread_mutex_unlock(&lock);
    }
    return NULL;
}
void* subber(){
    for(int x = 0; x < 1000009; x++){
        pthread_mutex_lock(&lock);
        i--;
        pthread_mutex_unlock(&lock);
    }
    return NULL;
}


int main(){
    pthread_t adder_thr;
    pthread_t subber_thr;
    pthread_create(&adder_thr, NULL, adder, NULL);
    pthread_create(&subber_thr, NULL, subber, NULL);
    pthread_join(adder_thr, NULL);
    pthread_join(subber_thr, NULL);
    printf("Done: %i\n", i);
    return 0;
}
