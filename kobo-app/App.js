import React from 'react';
import {StyleSheet, View, TouchableOpacity, Text } from 'react-native';

const sendRequest = (url) => {

  console.log("pressed!")

  // create a new XMLHttpRequest
  var xhr = new XMLHttpRequest()

  xhr.onreadystatechange = (e) => {
    if (xhr.readyState !== 4) {
      return;
    }
  
    if (xhr.status === 200) {
      console.log('success', xhr.responseText);
    } else {
      console.warn('error');
    }
  }
  // open the request with the verb and the url
  xhr.open('GET', url)
  xhr.send();
}

const App = () => {
  return (
    <View style={styles.container}>
      <TouchableOpacity
        style={styles.button}
        onPress={() => {sendRequest("http://192.168.178.57/right")}}
      >
        <Text>Next</Text>
      </TouchableOpacity>

      <TouchableOpacity
        style={styles.button}
        onPress={() => {sendRequest("http://192.168.178.57/left")}}
      >
        <Text>Back</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    paddingHorizontal: 10
  },
  button: {
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: "#869edb",
    padding: 10,
    height: 200,
    margin: 25
  },
});

export default App;