import { exec } from 'child_process';
import psList from 'ps-list';
import express from 'express';

const OBS_PATH = '/usr/bin/obs'; // Ensure this is the correct path
import cors from 'cors';


const checkOBS = async () => {
  try {
    console.log('Checking if OBS is running...');
    const processes = await psList();
    const isRunning = processes.some(process => process.name.includes('obs'));
    console.log(`OBS running: ${isRunning}`);
    return isRunning;
  } catch (error) {
    console.error('Error checking OBS processes:', error.message);
    return false;
  }
};

const startOBS = () => {
  console.log('Attempting to start OBS Studio...');
  exec(OBS_PATH, (error, stdout, stderr) => {
    if (error) {
      console.error(`Error starting OBS: ${error.message}`);
      return;
    }
    if (stderr) {
      console.error(`OBS stderr: ${stderr}`);
      return;
    }
    console.log(`OBS stdout: ${stdout}`);
  });
};

const ensureOBSRunning = async (req, res) => {
  try {
    const obsRunning = await checkOBS();
    if (!obsRunning) {
      console.log('OBS Studio is not running, attempting to start it...');
      startOBS();
      // Wait a bit for OBS to start before returning success
      setTimeout(() => {
        console.log('Returning success response after waiting for OBS to start...');
        res.status(200).json({ success: true });
      }, 5000);
    } else {
      console.log('OBS Studio is already running.');
      res.status(200).json({ success: true });
    }
  } catch (error) {
    console.error(`Failed to ensure OBS is running: ${error.message}`);
    res.status(500).json({ error: 'Internal server error' });
  }
};

const app = express();
app.use(cors());

app.get('/start-obs', ensureOBSRunning);

app.listen(3001, () => console.log('OBS controller backend running on port 3001'));
