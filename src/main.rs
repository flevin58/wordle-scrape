use anyhow::Result;
use scraper::{Html, Selector};
use reqwest::header::{HeaderMap, USER_AGENT};
use std::fs::write;

const URL: &str = "https://www.wordunscrambler.net/word-list/wordle-word-list";
const USER_AGENT_STR: &str = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3";
const WORDS_FILE: &str = "words.txt";

fn main() -> Result<()> {
    let html = fetch_html(URL)?;
    let words = extract_words(&html);
    write(WORDS_FILE, words.join("\n"))?;
    Ok(())
}

// Fetches page HTML over HTTP
fn fetch_html(url: &str) -> Result<String> {
    // Some sites block requests with no/odd User-Agent.
    // This is a simple, polite one.let mut headers = HeaderMap::new();
    let mut headers = HeaderMap::new();
    headers.insert(USER_AGENT, USER_AGENT_STR.parse().unwrap());

    // Build the Client
    let client = reqwest::blocking::Client::builder()
        .default_headers(headers)
        .build()?;

    let html = client
        .get(url)
        .send()?
        .text()?;

    Ok(html)
}

// Extracts word ranges from page HTML
fn extract_words(html: &str) -> Vec<String> {
    let mut words = Vec::<String>::new();
    let document = Html::parse_document(html);
    let a_selector = Selector::parse("h3.list-header + ul li a").unwrap();
    for a in document.select(&a_selector) {
        let text = a.text().collect::<String>();
        words.push(text);
    }
    words
}

#[cfg(test)]
mod tests {
    use crate::{extract_words, fetch_html};

    macro_rules! trim_lines {
        ($id:ident) => {
            $id
                . lines()
                .map(|line| line.trim())
                .collect::<String>()
        };
    }

    #[test]
    fn test_fetch_html() {
        let url = "https://jsonplaceholder.typicode.com/todos/1";
        let expected = r#"
            {
                "userId": 1,
                "id": 1,
                "title": "delectus aut autem",
                "completed": false
            }
        "#;
 
        let html = fetch_html(url).unwrap();

        assert_eq!(trim_lines!(html), trim_lines!(expected));
    }

    #[test]
    fn test_extract_words() {
        let html = r#"
            <div class="light-box text-left">
                <h3 class="list-header">Starting With A <span style="color: #94a3b8; font-weight: 400;">(143)</span></h3>
                <ul class="list-unstyled">
                    <li class="invert light">
                        <a href="/unscramble/aback">aback</a>
                    </li>
                    <li class="invert light">
                        <a href="/unscramble/abase">abase</a>
                    </li>
                </ul>
            </div>
        "#;

        let words = extract_words(html);
        assert_eq!(words, ["aback", "abase"]);
    }
}
