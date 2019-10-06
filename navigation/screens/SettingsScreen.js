import React from 'react';
import { ScrollView, View, Button, StyleSheet, Picker, Text, TextInput } from 'react-native';
import { Table, TableWrapper, Row, Rows, Col, Cols, Cell } from 'react-native-table-component';

export default class SettingsScreen extends React.Component {
  constructor(props){
    super(props);
    this.state = {
      tableHead: ['', 'Me', 'Average', 'Best'],
      tableData: [
        ['Trips To Date', '23', '13', '43'],
        ['Co2 Saved', '23232 ppm', '13131 ppm', '43443 ppm'],
        ['Tax Rebates', '$234', '$102', '$835'],
      ]
    }
  }

  render(){
    const state = this.state;
    return (
      <View style={styles.container}>
        <Table borderStyle={{borderWidth: 2, borderColor: '#c8e1ff'}}>
          <Row data={state.tableHead} style={styles.head} textStyle={styles.text}/>
          <Rows data={state.tableData} textStyle={styles.text}/>
        </Table>
      </View>
    );
  }
}

SettingsScreen.navigationOptions = {
  title: 'Superman\'s Scorecard',
};

const styles = StyleSheet.create({
  head: { height: 40, backgroundColor: '#f1f8ff' },
  text: { margin: 6 },
  container: {
    flex: 1,
    paddingTop: 15,
    backgroundColor: '#fff',
  },
});
