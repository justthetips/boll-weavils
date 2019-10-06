import React from 'react';
import { ScrollView, View, Button, StyleSheet, Picker, Text, TextInput } from 'react-native';
import GenerateForm from 'react-native-form-builder';
import { AppRegistry } from 'react-native';

export default class LinksScreen extends React.Component {
 
  constructor(props){
    super(props);
    this.state = {
      date: new Date('2020-06-12T14:42:42'),
      show: false,
    }
  }

  setDate = (event, date) => {
    date = date || this.state.date;

    this.setState({
      show: Platform.OS === 'ios' ? true : false,
      date,
    });
  }

  save = () => {
    const data = {};
    fetch('https://xajgj7zaxg.execute-api.us-east-1.amazonaws.com/prod/routes?route_number=1', {
         method: 'POST',
         body: JSON.stringify(data)
      })
      .then((response) => response.json())
      .then((responseJson) => {
        console.log(responseJson);
         
      })
      .catch((error) => {
         console.error(error);
      });

  }

  render(){
    const { show, date, mode } = this.state;

    return (
      <View style={styles.container}>
        
        <Text>Where do you want to go?</Text>
        <TextInput
        style={{ height: 40, borderColor: 'gray', borderWidth: 1 }}
        onChangeText={text => onChangeText(text)}
        
      />
        <Text>What time do you want to get there?</Text>
        <View style={{flex: 1, flexDirection: 'row'}}>
        <Picker
          selectedValue={this.state.hour}
          style={{height: 50, width: 100}}
          onValueChange={(itemValue, itemIndex) =>
            this.setState({hour: itemValue})
          }>
           
            <Picker.Item label="12" value="12" />
            <Picker.Item label="1" value="1" />
            <Picker.Item label="2" value="2" />
            <Picker.Item label="3" value="3" />
            <Picker.Item label="4" value="4" />
            <Picker.Item label="5" value="5" />
            <Picker.Item label="6" value="6" />
            <Picker.Item label="7" value="7" />
            <Picker.Item label="8" value="8" />
            <Picker.Item label="9" value="9" />
            <Picker.Item label="10" value="10" />
            <Picker.Item label="11" value="11" />
        </Picker>
        <Picker
          selectedValue={this.state.minute}
          style={{height: 50, width: 100}}
          onValueChange={(itemValue, itemIndex) =>
            this.setState({minute: itemValue})
          }>
            <Picker.Item label="00" value="00" />
            <Picker.Item label="15" value="15" />
            <Picker.Item label="30" value="30" />
            <Picker.Item label="45" value="45" />
        </Picker>
        <Picker
          selectedValue={this.state.ampm}
          style={{height: 50, width: 100}}
          onValueChange={(itemValue, itemIndex) =>
            this.setState({ampm: itemValue})
          }>
            <Picker.Item label="AM" value="AM" />
            <Picker.Item label="PM" value="PM" />
        </Picker>
        </View>
        <Button title="Save" onPress={this.save}></Button>
     
      </View>
    );
  }
}

LinksScreen.navigationOptions = {
  title: 'Sign Me Up',
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingTop: 15,
    backgroundColor: '#fff',
  },
});
