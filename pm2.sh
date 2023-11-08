#!/bin/bash
pm2 start pm2.json && pm2 save && pm2 log bec-webhooks
