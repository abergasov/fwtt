import axios from 'axios';
import crypto from 'crypto';

(async () => {
    console.log("starting client")
    try {
        console.log("requesting challenge without pow")
        const res = await axios.get('http://127.0.0.1:8000/quote');
    } catch (error) {
        console.error("got expected error:", error?.response?.data);
    }
    console.log("requesting challenge")
    const res = await axios.get('http://127.0.0.1:8000/challenges?num=1');
    console.log("got response:", res.data);

    let nonce = 0;
    let hash = "";
    if (res.data.algorithm === "sha256") {
        console.log("solving SHA256 pow")
        const {nonce: n, hash: h} = getSolutionSHA256(res.data.challenges[0], +res.data.difficulty);
        nonce = n;
        hash = h;
    } else if (res.data.algorithm === "scrypt") {
        console.log("solving SCRYPT pow")
        const {nonce: n, hash: h} = getSolutionScrypt(res.data.challenges[0], +res.data.difficulty, res.data.algo_params);
        nonce = n;
        hash = h;
    } else {
        throw new Error("unknown algorithm:", res.data.algorithm);
    }
    console.log(`sending solution: ${hash}, nonce: ${nonce}`);
    const config = {
        headers: {
            'X-Nonce': nonce,
            'X-Solution': hash,
            'X-Challenge': res.data.challenges[0],
        },
    }

    for (let i = 0; i < 10; i++) {
        console.log("requesting quote, iteration:", i)
        try {
            const quote = await axios.get('http://127.0.0.1:8000/quote', config)
            console.log("got response:", quote.data);
        } catch (error) {
            console.error("got expected error:", error?.response?.data);
            break;
        }
    }
})();


function getSolutionSHA256(challenge, diff) {
    let nonce = 0;
    while (true) {
        const hash = crypto.createHash('sha256');
        hash.update(`${challenge}${nonce}`);
        const str = hash.digest('hex');
        if (str.startsWith('0'.repeat(diff))) {
            return {nonce: nonce, hash: str};
        }
        nonce++;
    }
}

function getSolutionScrypt(challenge, diff, params) {
    let nonce = 0;
    while (true) {
        const buffer = new ArrayBuffer(4); // create an ArrayBuffer with 4 bytes
        const view = new DataView(buffer);
        view.setInt32(0, nonce, true); // set the number in the buffer, using little-endian byte order

        const solution = crypto.scryptSync(challenge, view, params.key_len, {
            N: params.n,
            r: params.r,
            p: params.p,
        })
        const str = solution.toString('hex');
        if (str.startsWith('0'.repeat(diff))) {
            return {nonce: nonce, hash: str};
        }
        nonce++;
    }
}
